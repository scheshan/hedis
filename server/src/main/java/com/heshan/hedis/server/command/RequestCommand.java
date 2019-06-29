package com.heshan.hedis.server.command;

import com.heshan.hedis.server.session.HedisSession;
import com.heshan.hedis.shared.codec.HedisMessage;

public class RequestCommand {

    private final HedisSession session;

    private final HedisMessage message;

    public RequestCommand(HedisSession session, HedisMessage message) {
        this.session = session;
        this.message = message;
    }

    public HedisSession session() {
        return session;
    }

    public HedisMessage message() {
        return message;
    }
}
