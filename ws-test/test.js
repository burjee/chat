function randomName() {
    let chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    let name = "";
    for (let i = 0; i < 12; i += 1) {
        let n = Math.floor(Math.random() * chars.length);
        name += chars[n];
    }
    return name
}

function join(userContext, _events, done) {
    const joinData = {
        method: "JOIN",
        name: randomName(),
        nonce: crypto.randomUUID()
    };
    userContext.vars.joinData = joinData;
    return done();
}

function message(userContext, _events, done) {
    const messageData = {
        method: "MESSAGE",
        message: {
            room: "room" + Date.now() % 2,
            text: "Test message",
        },
        nonce: crypto.randomUUID()
    };
    userContext.vars.messageData = messageData;
    return done();
}

module.exports = { join, message };
