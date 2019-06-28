package com.heshan.hedis.server.session;

import io.netty.channel.Channel;
import io.netty.util.AttributeKey;

import java.util.HashMap;
import java.util.Map;

/**
 * SessionManager
 *
 * @author heshan
 * @date 2019-06-28
 */
public class SessionManager {

    private Map<String, HedisSession> sessionMap;

    private static final AttributeKey<HedisSession> CHANNEL_SESSION_KEY = AttributeKey.valueOf("HedisSession");

    private static SessionManager instance = new SessionManager();

    private SessionManager() {
        sessionMap = new HashMap<>();
    }

    public void sessionInit(Channel channel) {
        HedisSession session = new HedisSession(channel);
        channel.attr(CHANNEL_SESSION_KEY).set(session);
        sessionMap.put(session.id(), session);

        System.out.printf("Session: %s init\r\n", session.id());
    }

    public void sessionClose(Channel channel) {
        HedisSession session = get(channel);
        sessionMap.remove(session.id());

        System.out.printf("Session: %s destroy\r\n", session.id());
    }

    public static SessionManager getInstance() {
        return instance;
    }

    public HedisSession get(String id) {
        return sessionMap.get(id);
    }

    public HedisSession get(Channel channel){
        return channel.attr(CHANNEL_SESSION_KEY).get();
    }
}
