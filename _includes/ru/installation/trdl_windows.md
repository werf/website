Убедитесь, что [Git](https://git-scm.com/download/win) версии 2.18.0 или новее, gpg и [Docker](https://docs.docker.com/get-docker) установлены. Дальнейшие инструкции должны выполняться в PowerShell.

{% offtopic title="Предоставьте вашему пользователю права на создание символьных ссылок" %}
> Если не хотите выполнять нижеуказанный PowerShell-скрипт, то то же самое вы можете сделать через GUI.

Выполните с правами администратора:
```powershell
$ntprincipal = new-object System.Security.Principal.NTAccount "$env:UserName"
$sidstr = $ntprincipal.Translate([System.Security.Principal.SecurityIdentifier]).Value.ToString()
$tmp = [System.IO.Path]::GetTempFileName()
secedit.exe /export /cfg "$($tmp)"
$currentSetting = ""
foreach($s in (Get-Content -Path $tmp)) {
    if ($s -like "SECreateSymbolicLinkPrivilege*") {
        $x = $s.split("=",[System.StringSplitOptions]::RemoveEmptyEntries)
        $currentSetting = $x[1].Trim()
    }
}
if ($currentSetting -notlike "*$($sidstr)*") {
    if ([string]::IsNullOrEmpty($currentSetting)) {
        $currentSetting = "*$($sidstr)"
    } else {
        $currentSetting = "*$($sidstr),$($currentSetting)"
    }
    $tmp2 = [System.IO.Path]::GetTempFileName()
    @"
[Unicode]
Unicode=yes
[Version]
signature="`$CHICAGO`$"
Revision=1
[Privilege Rights]
SECreateSymbolicLinkPrivilege = $($currentSetting)
"@ | Set-Content -Path $tmp2 -Encoding Unicode -Force
    cd (Split-Path $tmp2)
    secedit.exe /configure /db "secedit.sdb" /cfg "$($tmp2)" /areas USER_RIGHTS
}
```

Теперь выйдите из своей учётной записи Windows, затем зайдите в неё обратно и выполните следующую команду:
```powershell
gpupdate /force
```
{% endofftopic %}

[Установите trdl](https://github.com/werf/trdl/releases/) в `<диск>:\Users\<имя пользователя>\bin\trdl`. `trdl` будет отвечать за установку и обновление `werf`. Добавьте `<диск>:\Users\<имя пользователя>\bin\` в переменную окружения $PATH.

Добавьте `werf`-репозиторий в `trdl`:
```powershell
trdl add werf https://tuf.werf.io 1 b7ff6bcbe598e072a86d595a3621924c8612c7e6dc6a82e919abe89707d7e3f468e616b5635630680dd1e98fc362ae5051728406700e6274c5ed1ad92bea52a2
```
 
Для использования `werf` на рабочей машине мы рекомендуем настроить для `werf` _автоматическую активацию_. Для этого команда активации должна запускаться для каждой новой PowerShell-сессии. В PowerShell для этого обычно надо добавить команду активации в [$PROFILE-файл](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles). Команда активации `werf` для текущей PowerShell-сессии:
```powershell
. $(trdl use werf {{ include.version }} {{ include.channel }})
```

Для использования `werf` в CI вместо автоматической активации предпочитайте активацию `werf` вручную. Для этого выполните команду активации в начале вашей CI job, до вызова самого `werf`.

После активации `werf` должен быть доступен в той же PowerShell-сессии, в которой он был активирован:
```powershell
werf version
```
