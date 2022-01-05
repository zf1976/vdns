package parameter

import (
	"math/rand"
	"net/url"
	"strconv"
	"time"
	"vdns/lib/api/errs"
	"vdns/lib/api/models"
	"vdns/lib/auth"
	"vdns/lib/sign/compose"
	"vdns/lib/standard"
	"vdns/lib/standard/msg"
	"vdns/lib/standard/record"
	"vdns/vutil/strs"
	"vdns/vutil/vhttp"
)

type DnspodParameterProvider struct {
	credential        auth.Credential
	signatureComposer compose.SignatureComposer
	version           *standard.Standard
}

func NewDnspodParameterProvider(credential auth.Credential, signatureComposer compose.SignatureComposer) ParamaterProvider {
	return &DnspodParameterProvider{
		credential:        credential,
		signatureComposer: signatureComposer,
		version:           standard.DNSPOD_API_VERSION.String(),
	}
}

func (_this *DnspodParameterProvider) LoadDescribeParamater(request *models.DescribeDomainRecordsRequest, action *string) (*url.Values, error) {
	if request == nil {
		return nil, errs.NewApiError(msg.DESCRIBE_REQUEST_NOT_NIL)
	}
	// assert domain
	extractDomain, err := vhttp.ExtractDomain(strs.StringValue(request.Domain))
	if err != nil {
		return nil, errs.NewApiErrorFromError(err)
	}
	domain := extractDomain[0]
	rr := extractDomain[1]
	paramter := _this.loadCommonParamter(action)
	paramter.Set(DNSPOD_PARAMETER_DOMAIN, domain)

	// assert record type
	if !record.Support(request.RecordType) {
		return nil, errs.NewApiError(msg.RECORD_TYPE_NOT_SUPPORT)
	}
	paramter.Set(DNSPOD_PARAMETER_RECORD_TYPE, request.RecordType.String())

	// assert page size
	if request.PageSize != nil {
		paramter.Set(DNSPOD_PARAMETER_LIMIT, strconv.FormatInt(*request.PageSize, 10))
	}

	// assert offset start from 0
	if request.PageNumber != nil {
		paramter.Set(DNSPOD_PARAMETER_OFFSET, strconv.FormatInt(*request.PageNumber-1, 10))
	}

	// search and parse records by keyword, currently supports searching for host headers and record values
	if request.ValueKeyWord != nil {
		paramter.Set(DNSPOD_PARAMETER_KEY_WORD, *request.ValueKeyWord)
	}

	// assert rr key word
	if request.RRKeyWord != nil {
		paramter.Set(DNSPOD_PARAMETER_SUBDOMAIN_1, *request.RRKeyWord)
	} else if strs.NotEmpty(rr) {
		paramter.Set(DNSPOD_PARAMETER_SUBDOMAIN_1, rr)
	}
	return paramter, nil
}

func (_this *DnspodParameterProvider) LoadCreateParamater(request *models.CreateDomainRecordRequest, action *string) (*url.Values, error) {
	if request == nil {
		return nil, errs.NewApiError(msg.DESCRIBE_REQUEST_NOT_NIL)
	}

	// assert record type assert
	if request.RecordType != nil && !record.Support(*request.RecordType) {
		return nil, errs.NewApiError(msg.RECORD_TYPE_NOT_SUPPORT)
	}

	// assert value
	if request.Value == nil {
		return nil, errs.NewApiError(msg.RECORD_VALUE_NOT_SUPPORT)
	}

	// assert domain
	extractDomain, err := vhttp.ExtractDomain(strs.StringValue(request.Domain))
	if err != nil {
		return nil, errs.NewApiErrorFromError(err)
	}
	domain := extractDomain[0]
	rr := extractDomain[1]
	paramter := _this.loadCommonParamter(action)
	paramter.Set(DNSPOD_PARAMETER_DOMAIN, domain)
	paramter.Set(DNSPOD_PARAMETER_RECORD_TYPE, request.RecordType.String())
	paramter.Set(DNSPOD_PARAMETER_RECORD_LINE, DNSPOD_PARAMETER_DEFAULT)
	// assert rr
	if strs.IsEmpty(rr) {
		paramter.Set(DNSPOD_PARAMETER_SUBDOMAIN_2, record.PAN_ANALYSIS_RR_KEY_WORD.String())
	} else {
		paramter.Set(DNSPOD_PARAMETER_SUBDOMAIN_2, rr)
	}
	return paramter, nil
}

func (_this *DnspodParameterProvider) LoadUpdateParamater(request *models.UpdateDomainRecordRequest, action *string) (*url.Values, error) {
	if request == nil {
		return nil, errs.NewApiError(msg.DESCRIBE_REQUEST_NOT_NIL)
	}

	// assert record id
	if request.ID == nil {
		return nil, errs.NewApiError(msg.RECORD_ID_NOT_SUPPORT)
	}

	// assert record type assert
	if !record.Support(request.RecordType) {
		return nil, errs.NewApiError(msg.RECORD_TYPE_NOT_SUPPORT)
	}

	// assert value
	if request.Value == nil {
		return nil, errs.NewApiError(msg.RECORD_VALUE_NOT_SUPPORT)
	}

	// assert domain
	extractDomain, err := vhttp.ExtractDomain(strs.StringValue(request.Domain))
	if err != nil {
		return nil, errs.NewApiErrorFromError(err)
	}
	domain := extractDomain[0]
	rr := extractDomain[1]
	paramter := _this.loadCommonParamter(action)
	paramter.Set(DNSPOD_PARAMETER_RECORD_ID, *request.ID)
	paramter.Set(DNSPOD_PARAMETER_DOMAIN, domain)
	paramter.Set(DNSPOD_PARAMETER_RECORD_TYPE, request.RecordType.String())
	paramter.Set(DNSPOD_PARAMETER_RECORD_LINE, DNSPOD_PARAMETER_DEFAULT)

	// assert rr
	if strs.IsEmpty(rr) {
		paramter.Set(DNSPOD_PARAMETER_SUBDOMAIN_2, record.PAN_ANALYSIS_RR_KEY_WORD.String())
	} else {
		paramter.Set(DNSPOD_PARAMETER_SUBDOMAIN_2, rr)
	}

	return paramter, nil
}

func (_this *DnspodParameterProvider) LoadDeleteParamater(request *models.DeleteDomainRecordRequest, action *string) (*url.Values, error) {
	if request == nil {
		return nil, errs.NewApiError(msg.DESCRIBE_REQUEST_NOT_NIL)
	}

	// assert record id
	if request.ID == nil {
		return nil, errs.NewApiError(msg.RECORD_ID_NOT_SUPPORT)
	}

	// assert domain
	extractDomain, err := vhttp.ExtractDomain(strs.StringValue(request.Domain))
	if err != nil {
		return nil, errs.NewApiErrorFromError(err)
	}
	paramter := _this.loadCommonParamter(action)
	paramter.Set(DNSPOD_PARAMETER_RECORD_ID, *request.ID)
	paramter.Set(DNSPOD_PARAMETER_DOMAIN, extractDomain[0])
	return paramter, nil
}

func (_this *DnspodParameterProvider) loadCommonParamter(action *string) *url.Values {
	paramater := make(url.Values, 10)
	nonce := strconv.FormatInt(rand.Int63()+time.Now().UnixMilli(), 10)
	timestamp := strconv.FormatInt(time.Now().UnixMilli()/1000, 10)
	paramater.Set(DNSPOD_PARAMETER_ACTION, strs.StringValue(action))
	paramater.Set(DNSPOD_PARAMETER_NONCE, nonce)
	paramater.Set(DNSPOD_PARAMETER_TIMESTAMP, timestamp)
	paramater.Set(DNSPOD_PARAMETER_SECRET_ID, _this.credential.GetSecretId())
	paramater.Set(DNSPOD_PARAMETER_SIGNATUREMETHOD, _this.signatureComposer.SignatureMethod())
	paramater.Set(DNSPOD_PARAMETER_VERSION, _this.version.StringValue())

	return &paramater
}
