# reversebilinear2

バイリニアで2倍に拡大しちゃった画像を元のサイズにきれいに戻すやつ

```
USAGE: reversebilinear2.exe <IN> [<OUT>]
```

- PNGのみ対応
- フォルダを指定すると中のファイル全部処理する
- [経緯はこの辺で](https://twitter.com/chigiri_vrc/status/1692041747315270135)

## 原理

![image](https://github.com/chigirits/reversebilinear2/assets/61717977/ab36752b-88e4-4812-9c13-0f2e477664a0)

バイリニア拡大時に、たとえば x は a,b,c,d のリニアな内分点として画素値が混合されるから
$$x = 0.75 (0.75 a + 0.25 b) + 0.25 (0.75 c + 0.25 d)$$

4つの点について連立方程式を立てると

$$
\begin{pmatrix}
x \\
y \\
z \\
w \\
\end{pmatrix} = \frac{1}{16} \begin{pmatrix}
9 & 3 & 3 & 1 \\
3 & 9 & 1 & 3 \\
3 & 1 & 9 & 3 \\
1 & 3 & 3 & 9 \\
\end{pmatrix} \begin{pmatrix}
a \\
b \\
c \\
d \\
\end{pmatrix}
$$

逆行列を計算して整理すると

$$
\begin{pmatrix}
a \\
b \\
c \\
d \\
\end{pmatrix} = \frac{1}{4} \begin{pmatrix}
9 & -3 & -3 & 1 \\
-3 & 9 & 1 & -3 \\
-3 & 1 & 9 & -3 \\
1 & -3 & -3 & 9 \\
\end{pmatrix} \begin{pmatrix}
x \\
y \\
z \\
w \\
\end{pmatrix}
$$

これにより拡大前の画素値が分かる

## バイリニアで2倍しちゃった動画 → 同サイズ高画質化の手順まとめ

1. `ffmpeg -i movie-in.mp4 -vcodec png sequence\image-%06d.png`
2. `sequence\*` を waifu2x-caffe でノイズ除去（レベル2・UpPhoto・TTAオフ・分割サイズ512）
3. その結果を reversebilinear2（本ツール）で半分サイズに復元、`sequence.half\*` に出力
4. `sequence.half\*` をUpscaylで拡大（REMACRIx4）、`sequence.remacri\*` に出力
5. `ffmpeg -r 60 -i sequence.remacri\image-%06d.png -sws_flags lanczos+accurate_rnd -s 3840x2160 -vcodec libx264 -pix_fmt yuv420p -b:v 90M movie-out.mp4`
