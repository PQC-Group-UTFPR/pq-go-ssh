param(
    [Parameter(Mandatory=$true, Position=0)]
    [string]$GB
)

[Console]::InputEncoding  = [System.Text.Encoding]::UTF8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding           = [System.Text.Encoding]::UTF8

$LibFile = "$PSScriptRoot\biblioteca_gb.txt"
if (-not (Test-Path $LibFile)) {
    New-Item -ItemType File -Path $LibFile | Out-Null
}

$ListaAtual = Get-Content $LibFile -Encoding UTF8 | Where-Object { $_.Trim() -ne "" }

if ($ListaAtual -notcontains $GB) {
    Add-Content -Path $LibFile -Value $GB -Encoding UTF8
    Write-Host "[BIBLIOTECA] Novo valor '$GB GB' cadastrado! Nas proximas vezes, digite apenas: $GB" -ForegroundColor Green
} else {
    Write-Host "[BIBLIOTECA] Executando '$GB GB' da biblioteca..." -ForegroundColor Cyan
}

# Cria o atalho direto com o próprio número (ex: digitar apenas "5" roda o teste de 5 GB)
Set-Item -Path "Function:global:gb$GB" -Value ([scriptblock]::Create("& '$PSScriptRoot\rodar_teste.ps1' '$GB'"))

Write-Host "`n[DOCKER] Executando benchmark de $GB GB..." -ForegroundColor Yellow
$Saida = docker run --rm -e TRANSFER_GB=$GB pq-ssh-bench
$Saida | ForEach-Object { Write-Host $_ }
$Saida | Out-File -FilePath "$PSScriptRoot\resultado_${GB}gb.txt" -Encoding UTF8
Write-Host "[OK] Resultado salvo em resultado_${GB}gb.txt" -ForegroundColor Green