# zero-gorm

zero-gorm is to adapt to go-zero with gorm.

# Usage
```
> go get github.com/yfanswer/zero-gorm
> replace your template/model with zero-gorm/template/1.2.4/model
> goctl model mysql ddl -src ./xxx.sql -dir ./ -c --home /xxx/template
```

# Remarks
- The version of your goctl should correspond to the version in the template directory.
- Need in the code with import github.com/yfanswer/zero-gorm/db.
