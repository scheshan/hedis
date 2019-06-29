package com.heshan.hedis.server.command;

import com.heshan.hedis.server.session.HedisSession;
import com.heshan.hedis.shared.codec.HedisMessage;

public class RequestWrapper {

    private final HedisMessage message;

    private final HedisSession session;

    public RequestWrapper(HedisMessage message, HedisSession session) {
        this.message = message;
        this.session = session;
    }

    public HedisMessage message() {
        return message;
    }

    public HedisSession session() {
        return session;
    }
}
