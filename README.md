# printer

Go 打印工具集：

- `printer`：Windows 专用（winspool）打印封装，可静默发送 RAW 数据到打印机（适合针式/ESC-P）。
- `xprint`：跨平台静默打印封装（Windows / macOS / Linux），支持 RAW / PDF / 网页(URL→PDF)。

## 安装

```bash
go get github.com/w6xian/printer
```

## 包说明

### printer（Windows）

`printer` 包封装了 Windows Spooler API（winspool.drv），用于：

- 获取默认打印机、枚举打印机
- 启动打印任务/页面并写入数据（RAW / XPS_PASS）

示例（打印一段文本）：

```go
package main

import (
	"fmt"
	"log"

	"github.com/w6xian/printer"
)

func main() {
	name, err := printer.Default()
	if err != nil {
		log.Fatal(err)
	}

	p, err := printer.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	if err := p.StartRawDocument("demo"); err != nil {
		log.Fatal(err)
	}
	defer p.EndDocument()

	if err := p.StartPage(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(p, "Hello %s\r\n", name)

	if err := p.EndPage(); err != nil {
		log.Fatal(err)
	}
}
```

针式打印机常见用法是直接写入 ESC/P 等控制指令（同样走 RAW 通道）。

### xprint（跨平台）

`xprint` 提供统一接口：

- `PrintRaw`：RAW/ESC-P 静默打印
- `PrintPDF`：PDF 静默打印
- `PrintURL`：网页静默打印（先用 headless Chrome/Chromium/Edge 转 PDF，再打印）

平台后端：

- Windows：RAW 走 `printer`（winspool），PDF 走 SumatraPDF 静默打印
- macOS/Linux：RAW/PDF 走 CUPS（`lp`/`lpstat`）

示例（静默打印 PDF）：

```go
package main

import (
	"log"

	"github.com/w6xian/printer/xprint"
)

func main() {
	err := xprint.PrintPDF("a.pdf", xprint.Options{
		Printer: "Your Printer Name",
		Copies:  1,
		JobName: "a.pdf",
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

示例（网页 URL 静默打印）：

```go
err := xprint.PrintURL("https://example.com", xprint.Options{
	Printer: "Your Printer Name",
	Copies:  1,
	JobName: "example",
})
```

可选环境变量：

- `XPRINT_CHROME`：指定 Chrome/Chromium/Edge 可执行文件路径（用于 URL→PDF）
- `XPRINT_SUMATRA`：Windows 下指定 `SumatraPDF.exe` 路径（用于 PDF 打印）

## 命令行工具

仓库内置 `cmd/print`，支持打印文本/PDF/URL：

```bash
go run ./cmd/print -l
go run ./cmd/print -p "Printer Name" ./demo.txt
go run ./cmd/print -p "Printer Name" ./a.pdf
go run ./cmd/print -p "Printer Name" https://example.com
```

## 注意事项

- PDF/网页打印依赖外部渲染/打印后端（Windows：SumatraPDF；macOS/Linux：CUPS；网页转 PDF：Chrome/Chromium/Edge）。
- 针式打印机建议优先用 RAW + 指令集方式（速度快、对齐稳定），避免整页位图渲染。
