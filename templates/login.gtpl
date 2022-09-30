<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Login</title>
    </head>
    <body>
        <h2>Login</h2>
        <hr>
        <form name="login_form" action="{{ .action }}" method="post">
            <input type="id" name="user_id" placeholder="User ID"><br>
            <input type="password" name="password"placeholder="Password">
            <button type="submit">Login</button>
        </form>
        {{ .message }}
    </body>
</html>