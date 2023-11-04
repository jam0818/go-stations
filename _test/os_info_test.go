package _test

import (
	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestOSInfo(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/healthz", nil)
	if err != nil {
		t.Error("リクエストの作成に失敗しました。", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("リクエストの送信に失敗しました。", err)
		return
	}
	want := "{\"message\":\"OK\"}\n"
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error("レスポンスの読み込みに失敗しました。", err)
		return
	}

	if string(got) != want {
		t.Errorf("期待していない内容です, got = %s, want = %s\n", got, want)
		return
	}
}
