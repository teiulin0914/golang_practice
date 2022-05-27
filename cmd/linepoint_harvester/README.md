# LINE POINTS Harvester
以爬蟲自動領取 LINE POINTS 的簡單腳本

## Build
建構執行檔 `linepoint_harvester` (windows 為 `linepoint_harvester.exe`)。
```
$ go build .
```

## Before running
1. 安裝 Chrome 瀏覽器
2. 至 https://points.line.me/pointcode/ 手動領取一次 LINE POINT，先行閱讀官方同意事項

## Run
1. 建立檔名為 `codes` 的文字檔，並填入欲領取的 LINE POINTS code (一行一筆)。
2. 將 `codes` 移至執行檔 `linepoint_harvester` 相同路徑下。
3. 執行 `linepoint_harvester` 。
    ```
    $ ./linepoint_harvester
    ```
4. 依序填入 LINE 帳號、密碼，Enter 後即會開始運行爬蟲。
5. 完成後，領取失敗的 code 會記錄在 `failed_codes` 中，發生非預期錯誤的 code 則會記錄在 `error_codes` 中。
