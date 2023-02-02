
function Service-Exist {
    param (
        $ServiceName
    )
    $ServiceObj = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($ServiceObj -eq $null) {
        return $false
    }
    return $true
}

function Run-Service {
    param (
        $ServiceName
    )

    if (!(Service-Exist $ServiceName)) {
        echo "Error: $ServiceName is not existed" | Out-Default
        return $false
    }

    $ServiceObj = Get-Service -Name $ServiceName
    if ($ServiceObj.Status -eq "Running") {
        echo "$($ServiceObj.Name) is running" | Out-Default
        return $true
    }

    Start-Service -Name $ServiceName
    $ServiceObj = Get-Service -Name $ServiceName
    if ($ServiceObj.Status -eq "Running") {
        echo "$($ServiceObj.Name) is running" | Out-Default
        return $true
    }
    return $false
}

echo "RUNNING SERVICE"
echo ""

$Services = @(
    "mysql80"
)

foreach ($servicename in $Services) {
    if (!(Run-Service $servicename)) {
        echo "Error: service is exist, but can't start service"
        exit
    }
}

echo ""
echo "ALL SERVICE IS RUNNING"
echo ""

fresh