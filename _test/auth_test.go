package _test

import (
	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Errorf("dbPathのセットに失敗しました。%v", err)
		return
	}

	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		t.Errorf("データベースの作成に失敗しました: %v", err)
		return
	}

	t.Cleanup(func() {
		if err := todoDB.Close(); err != nil {
			t.Errorf("データベースのクローズに失敗しました: %v", err)
			return
		}
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
			return
		}
	})
	r := router.NewRouter(todoDB)
	srv := httptest.NewServer(r)
	defer srv.Close()

	err = os.Setenv("BASIC_AUTH_USER_ID", "trueid")
	if err != nil {
		t.Errorf("IDの設定に失敗しました: %v", err)
	}
	err = os.Setenv("BASIC_AUTH_PASSWORD", "truepassword")
	if err != nil {
		t.Errorf("passwordの設定に失敗しました: %v", err)
	}

	testcases := map[string]struct {
		ID                 string
		Password           string
		WantHTTPStatusCode int
	}{
		"Different API": {
			ID:                 "trueid",
			Password:           "truepassword",
			WantHTTPStatusCode: http.StatusOK,
		},
		"Valid auth": {
			ID:                 "trueid",
			Password:           "truepassword",
			WantHTTPStatusCode: http.StatusOK,
		},
		"Invalid auth": {
			ID:                 "falseid",
			Password:           "falsepassword",
			WantHTTPStatusCode: http.StatusUnauthorized,
		},
		"Empty auth": {
			ID:                 "",
			Password:           "",
			WantHTTPStatusCode: http.StatusUnauthorized,
		},
		"Send no auth": {
			ID:                 "trueid",
			Password:           "truepassword",
			WantHTTPStatusCode: http.StatusUnauthorized,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			var req *http.Request
			switch name {
			case "Different API":
				req, err = http.NewRequest(http.MethodGet, srv.URL+"/healthz", nil)
				if err != nil {
					t.Errorf("リクエストの作成に失敗しました: %v", err)
					return
				}
				req.SetBasicAuth(tc.ID, tc.Password)
			case "Send no auth":
				req, err = http.NewRequest(http.MethodGet, srv.URL+"/os_info", nil)
				if err != nil {
					t.Errorf("リクエストの作成に失敗しました: %v", err)
					return
				}
			default:
				req, err = http.NewRequest(http.MethodGet, srv.URL+"/os_info", nil)
				if err != nil {
					t.Errorf("リクエストの作成に失敗しました: %v", err)
					return
				}
				req.SetBasicAuth(tc.ID, tc.Password)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("リクエストの送信に失敗しました: %v", err)
				return
			}
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("レスポンスのクローズに失敗しました: %v", err)
					return
				}
			})
			if resp.StatusCode != tc.WantHTTPStatusCode {
				t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, tc.WantHTTPStatusCode)
				return
			}
		})
	}
}
