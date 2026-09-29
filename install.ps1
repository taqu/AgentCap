$installerPath = Join-Path $env:TEMP 'codex-install.ps1'

$previousNonInteractive = [Environment]::GetEnvironmentVariable('CODEX_NON_INTERACTIVE', 'Process')
try {
    Invoke-WebRequest -Uri 'https://chatgpt.com/codex/install.ps1' -OutFile $installerPath -ErrorAction Stop
    $env:CODEX_NON_INTERACTIVE = '1'
    & 'C:\Program Files\PowerShell\7\pwsh.exe' `
        -NoProfile -ExecutionPolicy Bypass `
        -File $installerPath -Release latest

    if ($LASTEXITCODE -ne 0) {
        throw "Codexの更新に失敗しました: exit=$LASTEXITCODE"
    }

    codex --version
}
finally {
    [Environment]::SetEnvironmentVariable('CODEX_NON_INTERACTIVE', $previousNonInteractive, 'Process')
}

