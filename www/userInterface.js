function onSubmit(form) {

    var data = JSON.stringify($(form).serializeArray());

    console.log(data);

    // TODO: Send JSON to HTTP Server
}
