$LibFile = "$PSScriptRoot\biblioteca_gb.txt"
if (Test-Path $LibFile) {
    Get-Content $LibFile -Encoding UTF8 | Where-Object { $_.Trim() -ne "" } | ForEach-Object {
        $val = $_.Trim()
        $LibFile = "$PSScriptRoot\biblioteca_gb.txt"
if (Test-Path $LibFile) {
    Get-Content $LibFile -Encoding UTF8 | Where-Object { $_.Trim() -ne "" } | ForEach-Object {
        $val = $_.Trim()
        Set-Item -Path "Function:global:gb$val" -Value ([scriptblock]::Create("& '$PSScriptRoot\rodar_teste.ps1' '$val'"))
    }
}
    }
}