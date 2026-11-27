
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("run", "build", "test")]
    [string]$Action
)

# Переходим в папку, где расположен скрипт
Set-Location $PSScriptRoot

switch ($Action) {
    "run" {
        Write-Host "Starting Go server..."
        go run ./cmd/server
    }

    "build" {
        Write-Host "Building Go application..."

        if (-not (Test-Path "bin")) {
            New-Item -ItemType Directory -Path "bin" | Out-Null
        }

        go build -o ./bin/server.exe ./cmd/server
    }

    "test" {
        Write-Host "Running Go tests..."
        go test ./...
    }
}
