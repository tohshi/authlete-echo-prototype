<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Login</title>
    </head>
    <body>
        <h2>Consent</h2>
        <hr>
        {{ if .scopes }}
            <table>
            <tr>
                <th>scope</th>
                <th>description</th>
            </tr>
            {{ range .scopes }}
            <tr>
                <td>{{ .Name }}</td>
                <td>{{ .Description }}</td>
            </tr>
            {{ end }}
        {{ end }}
        </table>
        <form name="consent_form" action="{{ .action }}" method="post">
        <button type="submit" name="consent" value="agree">Agree</button>
        <button type="submit" name="consent" value="reject">Reject</button>
        </form>
    </body>
</html>