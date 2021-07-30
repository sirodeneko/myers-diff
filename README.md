# myers-diff
myers 算法的实现

# test
```
go run main.go a.txt b.txt
```
output
```
 func create() *Patches {
-       return &Patches{originals: make(map[reflect.Value][]byte), values: make(map[reflect.Value]reflect.Value), valueHolders: make(map[reflect.Value]reflect.Value)}
+       return &Patches{originals: make(map[reflect.Value][]byte), values: make(map[reflect.Value]reflect.Value)}
 }
 
 func NewPatches() *Patches {
        return create()
 }
 func (this *Patches) ApplyFunc(target, double interface{}) *Patches {
        t := reflect.ValueOf(target)
        d := reflect.ValueOf(double)
        return this.ApplyCore(t, d)
 }

```