package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type TestGame struct {
	elements tendon.Elements
}

func (g *TestGame) Update() error {
	// ★ 追加: Tendon全体のグローバル更新（キー入力監視など）を毎フレーム実行する
	tendon.GlobalUpdate()

	g.elements.Update(0, 0)
	return nil
}

func (g *TestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	g.elements.Draw(screen)
}

func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func TestInteractive(t *testing.T) {
	// デバッグモードをONにする
	tendon.DebugMode = true

	// 1. コマンドプロンプトの作成と設定
	// 画面下部に配置する (幅640, 高さ40)
	prompt := tendon.NewCommandPrompt(640, 40)
	prompt.XRelativeToParent = 0
	prompt.YRelativeToParent = 480 - 40
	
	// コマンドが確定（Enter）された時の挙動を定義
	prompt.OnExecute = func(command string, target *tendon.Element) {
		if target == nil {
			fmt.Printf("🚀 グローバルコマンド実行: %s\n", command)
		} else {
			fmt.Printf("Targeted 🎯 要素 [%s] へのコマンド実行: %s\n", target.Name, command)
			// ここで command の内容に応じて target.CenterInScreen(640, 480) などを呼ぶ処理を追加できる
		}
	}

	// 2. グローバルなキー入力を監視してプロンプトを開く
	tendon.OnKeyJustPressed = func(key ebiten.Key) {
		if key == ebiten.KeySlash {
			// まだ開いていなければ、グローバル(target=nil)としてプロンプトを表示
			if !prompt.Visible {
				prompt.Open(nil)
			}
		}
	}

	fmt.Println("=================================================")
	fmt.Println("コマンドプロンプト機能が有効です。")
	fmt.Println("【テスト方法】")
	fmt.Println("1. [/] キーを押すと、画面下にコマンドプロンプトが開きます。")
	fmt.Println("2. 文字を入力して Enter を押すとコンソールに内容が出ます。")
	fmt.Println("3. 青いパネルを【右クリック】すると、その要素専用のプロンプトが開きます。")
	fmt.Println("=================================================")

	// 3. 大きな「親パネル」
	panel := tendon.NewButton(50, 50, 300, 300, "Parent Panel", color.RGBA{80, 80, 150, 255})
	panel.Z = 1
	panel.Draggable = true
	panel.Name = "ParentPanel"
	
	// 初期状態で画面中央に配置する
	panel.CenterInScreen(640, 480)

	panel.OnMouseEnter = func(e *tendon.Element) {
		fmt.Println("🟦 親パネルにマウスが【入りました】")
	}
	panel.OnMouseLeave = func(e *tendon.Element) {
		fmt.Println("🟦 親パネルからマウスが【出ました】")
	}

	// ★ 追加: 右クリックでこの要素をターゲットにしたプロンプトを開く
	panel.OnRightClick = func(e *tendon.Element) {
		prompt.Open(e)
	}

	// 4. パネルの中に入れる「子ボタン」
	childBtn := tendon.NewButton(100, 150, 100, 50, "Child Btn", color.RGBA{200, 80, 80, 255})
	childBtn.Name = "ChildButton"

	childBtn.OnMouseEnter = func(e *tendon.Element) {
		fmt.Println("  🟥 子ボタンにマウスが【入りました】")
	}
	childBtn.OnMouseLeave = func(e *tendon.Element) {
		fmt.Println("  🟥 子ボタンからマウスが【出ました】")
	}
	childBtn.OnLeftClick = func(e *tendon.Element) {
		fmt.Println("【確認】子ボタンがクリックされました！")
	}

	panel.AppendChild(childBtn)

	// 5. もう一つの要素
	otherElem := tendon.NewButton(400, 50, 100, 100, "Other", color.RGBA{100, 100, 100, 255})
	otherElem.Z = 2
	otherElem.Name = "OtherElement"

	// 6. ゲームの実行（プロンプトのElementも忘れずに追加する）
	game := &TestGame{
		elements: tendon.Elements{panel, otherElem, prompt.Element},
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tendon Command Prompt Test")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}