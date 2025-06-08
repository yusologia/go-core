package logiaapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	logiares "github.com/yusologia/go-core/v2/response"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
)

type UploadFile struct {
	File        multipart.File
	FileHandler *multipart.FileHeader
}

type LogiaAPIOption struct {
	Headers map[string]string
}

type LogiaAPI interface {
	Get(url string, parameter url.Values) logiares.ResponseSuccessWithPagination
	Post(url string, payload any) logiares.ResponseSuccessWithPagination
	PostMultipart(url string, payload any) logiares.ResponseSuccessWithPagination
	Patch(url string, payload any) logiares.ResponseSuccessWithPagination
	Put(url string, payload any) logiares.ResponseSuccessWithPagination
	Delete(url string, payload any) logiares.ResponseSuccessWithPagination
}

func NewLogiaAPI(opt ...LogiaAPIOption) LogiaAPI {
	xa := logiaAPI{}
	if len(opt) > 0 {
		xa.option = opt[0]
	}

	return &xa
}

type logiaAPI struct {
	contentType string
	option      LogiaAPIOption
	payload     *bytes.Buffer
}

func (xa *logiaAPI) Get(url string, parameter url.Values) logiares.ResponseSuccessWithPagination {
	if filter := parameter.Encode(); filter != "" {
		url += "?" + filter
	}

	return xa.callAPI("GET", url)
}

func (xa *logiaAPI) Post(url string, payload any) logiares.ResponseSuccessWithPagination {
	xa.setJSONPayload(payload)

	return xa.callAPI("POST", url)
}

func (xa *logiaAPI) PostMultipart(url string, payload any) logiares.ResponseSuccessWithPagination {
	xa.setMultipartPayload(payload)

	return xa.callAPI("POST", url)
}

func (xa *logiaAPI) Patch(url string, payload any) logiares.ResponseSuccessWithPagination {
	xa.setJSONPayload(payload)

	return xa.callAPI("PATCH", url)
}

func (xa *logiaAPI) Put(url string, payload any) logiares.ResponseSuccessWithPagination {
	xa.setJSONPayload(payload)

	return xa.callAPI("PUT", url)
}

func (xa *logiaAPI) Delete(url string, payload any) logiares.ResponseSuccessWithPagination {
	xa.setJSONPayload(payload)

	return xa.callAPI("DELETE", url)
}

/** --- UNEXPORTED FUNCTIONS --- */

func (xa *logiaAPI) callAPI(method string, url string) logiares.ResponseSuccessWithPagination {
	if xa.contentType == "" {
		xa.contentType = "application/json"
	}

	payload := xa.payload
	if payload == nil {
		payload = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		logiares.ErrLogiaAPI(fmt.Sprintf("Set new request external api: %v", err.Error()))
	}
	req.Header.Set("Content-Type", xa.contentType)

	if len(xa.option.Headers) > 0 {
		for key, value := range xa.option.Headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logiares.ErrLogiaAPI(err.Error())
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logiares.ErrLogiaAPI(fmt.Sprintf("Read response external api: %v", err.Error()))
	}

	if resp.StatusCode > 299 {
		var errRes logiares.ResponseError
		err = json.Unmarshal(responseBody, &errRes)
		if err != nil {
			logiares.ErrLogiaAPI(fmt.Sprintf("Read error response: %v", err.Error()))
		}

		logiares.ErrLogiaAPI(fmt.Sprintf("%s. Internal-Msg: %s", errRes.Status.Message, errRes.Status.InternalMsg))
	}

	var success logiares.ResponseSuccessWithPagination
	err = json.Unmarshal(responseBody, &success)
	if err != nil {
		logiares.ErrLogiaAPI(fmt.Sprintf("Read success response: %v", err.Error()))
	}

	return success
}

func (xa *logiaAPI) setMultipartPayload(payload any) {
	writer := multipart.NewWriter(xa.payload)

	fields := make(map[string]string)
	xa.structToFormFields(payload, "", fields, writer)

	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	xa.contentType = writer.FormDataContentType()

	writer.Close()
}

func (xa *logiaAPI) setJSONPayload(payload any) {
	if payload != nil {
		payloadByte, err := json.Marshal(payload)
		if err != nil {
			logiares.ErrLogiaAPI(fmt.Sprintf("Marshal payload failed: %v", err.Error()))
		}

		xa.payload = bytes.NewBuffer(payloadByte)
	}
}

func (xa *logiaAPI) structToFormFields(v any, parent string, out map[string]string, writer *multipart.Writer) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("json")
		if tag == "" {
			continue
		}

		uploadFile, ok := fieldVal.Interface().(UploadFile)
		if ok && uploadFile.File != nil && uploadFile.FileHandler != nil {
			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, tag, filepath.Base(uploadFile.FileHandler.Filename)))
			header.Set("Content-Type", xa.getMimeType(uploadFile.File, uploadFile.FileHandler))

			part, err := writer.CreatePart(header)
			if err != nil {
				logiares.ErrLogiaAPI(fmt.Sprintf("Create multipart/form part is failed: %v", err.Error()))
			}

			if _, err = io.Copy(part, uploadFile.File); err != nil {
				logiares.ErrLogiaAPI(fmt.Sprintf("Copy file to multipart/form part is failed: %v", err.Error()))
			}

			continue
		}

		key := tag
		if parent != "" {
			key = parent + "[" + tag + "]"
		}

		switch fieldVal.Kind() {
		case reflect.Struct:
			xa.structToFormFields(fieldVal.Interface(), key, out, writer)
		default:
			out[key] = fmt.Sprint(fieldVal.Interface())
		}
	}
}

func (xa *logiaAPI) getMimeType(file multipart.File, handler *multipart.FileHeader) string {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		log.Panicf("Unable to reading file: %v", err)
	}

	mimeTypeSystem := http.DetectContentType(buf[:n])
	if mimeTypeSystem == "application/zip" {
		ext := strings.ToLower(filepath.Ext(handler.Filename))
		switch ext {
		case ".xlsx":
			return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case ".docx":
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		}
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Panicf("Unable to reset file pointer: %v", err)
	}

	return mimeTypeSystem
}
