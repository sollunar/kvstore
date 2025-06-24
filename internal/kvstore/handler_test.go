package kvstore

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sollunar/kvstore-api/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func bootstrap(t *testing.T) (*Handler, *mock.MockStorage, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockStore := mock.NewMockStorage(ctrl)
	kvService := NewKVService(mockStore)
	h := &Handler{kvservice: kvService}
	return h, mockStore, ctrl
}

func TestGetHandler(t *testing.T) {
	randKey := randomString(5)
	randValue := randomString(5)

	testCases := []struct {
		name          string
		key           string
		buildStubs    func(mockStore *mock.MockStorage)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			key:  randKey,
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().
					Get(randKey).
					Times(1).
					Return(randValue, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var body SetRequest
				err := json.NewDecoder(recorder.Body).Decode(&body)
				require.NoError(t, err)
				require.Equal(t, randKey, body.Key)
				require.Equal(t, randValue, body.Value)
			},
		},
		{
			name: "NOT_FOUND",
			key:  "missing",
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().
					Get("missing").
					Times(1).
					Return("", ErrKeyNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "BAD_REQUEST_EMPTY_KEY",
			key:        "",
			buildStubs: func(store *mock.MockStorage) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, mockStore, ctrl := bootstrap(t)
			defer ctrl.Finish()

			tc.buildStubs(mockStore)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get?key=%s", tc.key), nil)
			rec := httptest.NewRecorder()

			h.Get(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func TestSetHandler(t *testing.T) {
	randKey := randomString(5)
	randValue := randomString(5)

	testCases := []struct {
		name          string
		body          string
		buildStubs    func(mockStore *mock.MockStorage)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: fmt.Sprintf(`{"key": "%s", "value": "%s"}`, randKey, randValue),
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().
					Set(randKey, randValue).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, rec.Code)
			},
		},
		{
			name: "BAD_REQUEST_INVALID_JSON",
			body: `invalid-json`,
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().Set(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "BAD_REQUEST_KEY_OR_VALUE",
			body: fmt.Sprintf(`{"key": "", "value": "%s"}`, randValue),
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().Set(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "INTERNAL_ERROR",
			body: fmt.Sprintf(`{"key": "%s", "value": "%s"}`, randKey, randValue),
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().
					Set(randKey, randValue).
					Times(1).
					Return(ErrInternal)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, mockStore, ctrl := bootstrap(t)
			defer ctrl.Finish()

			tc.buildStubs(mockStore)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/set", strings.NewReader(tc.body))
			require.NoError(t, err)

			h.Set(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func TestDeleteKeyRequest(t *testing.T) {
	randKey := randomString(5)

	testCases := []struct {
		name          string
		key           string
		buildStubs    func(mockStore *mock.MockStorage)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			key:  randKey,
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().
					Delete(randKey).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name: "NOT_FOUND",
			key:  randKey,
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().
					Delete(randKey).
					Times(1).
					Return(ErrKeyNotFound)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "BAD_REQUEST_EMPTY_KEY",
			key:  "",
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().Delete(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "INTERNAL_ERROR",
			key:  randKey,
			buildStubs: func(store *mock.MockStorage) {
				store.EXPECT().
					Delete(randKey).
					Times(1).
					Return(ErrInternal)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, mockStore, ctrl := bootstrap(t)
			defer ctrl.Finish()

			tc.buildStubs(mockStore)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/delete?key=%s", tc.key), nil)

			h.Delete(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[seededRand.Intn(len(letterBytes))]
	}
	return string(b)
}
