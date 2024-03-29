# go-julia
julia set renderer powered by tcell. default constant $c$ is -0.8+0.156i

calculated via $f(z)=z^{2}+c$

### usage
```
Usage: go-julia [OPTION]... [REAL] [IMAGINARY]
Render Julia set at given complex constant.
Arrow keys to move, +/- to zoom.
Hold r or s to move between constants.

REAL and IMAGINARY are the real and imaginary components of a complex constant.
With no REAL and IMAGINARY, constant is -0.8+0.156i.

  -c1 string
        color 1 in RGB list format (default "0,255,0")
  -c2 string
        color 2 in RGB list format (default "0,0,0")
  -it int
        maximum iterations (default 100)
```
