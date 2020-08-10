[Home](https://emicklei.github.io/melrose)

## Example: iterator

```javascript
bpm(90)
s = sequence('c e g b')
i = iterator('(1 2) (3 4)', '(1 3) (2 4)', '1 (2 3) 4', '(1 4) (2 3)')
m = sequencemap(i,s)
l = loop(m,next(i))
```