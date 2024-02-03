package main

import (
	"log"
	"os"

	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv("API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// 新しいWord文書を作成
	doc := document.New()

	// 文書に段落を追加
	para := doc.AddParagraph()

	// 段落にテキストを追加
	run := para.AddRun()
	run.AddText("Hello, World!")

	// 文書をファイルに保存
	err := doc.SaveToFile("example.docx")
	if err != nil {
		log.Fatal(err)
	}
}
