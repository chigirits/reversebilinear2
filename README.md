# reversebilinear2

バイリニアで2倍に拡大しちゃった画像を元のサイズにきれいに戻すやつ

```
USAGE: reversebilinear2.exe <IN> [<OUT>]
```

- PNGのみ対応
- フォルダを指定すると中のファイル全部処理する
- [経緯はこの辺で](https://twitter.com/chigiri_vrc/status/1692041747315270135)

## メモ

![image](https://github.com/chigirits/reversebilinear2/assets/61717977/ab36752b-88e4-4812-9c13-0f2e477664a0)

バイリニア拡大時に、たとえば x は a,b,c,d のリニアな内分点として画素値が混合されるから
$$x = 0.75 (0.75 a + 0.25 b) + 0.25 (0.75 x + 0.25 d)$$

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

