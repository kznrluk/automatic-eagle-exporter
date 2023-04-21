# sdweb-eaglepack

[AUTOMATIC1111/stable-diffusion-webui](https://github.com/AUTOMATIC1111/stable-diffusion-webui) で作成された画像を、画像管理ソフトEagleにPNGInfoの情報とともに送信するためのソフトウェアです。

指定されたディレクトリに格納されている画像をEagleの素材パックである `.eaglepack` に変換します。ダブルクリックやD&Dで簡単にファイルをライブラリに追加することができます。

## Usage
対象のディレクトリやファイルをGLOBパターンを用いて指定します。

```
> sdweb-eaglepack ./txt2img-images/**/*.png
```

終わるとカレントディレクトリに `.eaglepack` が作成されます。

## 既知の問題
- 画像の色を抽出したパレットが生成されない
- サムネイルが生成されない

上記の問題はサムネイルの再生成を行うことで解消できます。