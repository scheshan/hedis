package com.heshan.hedis.server.session;

import com.heshan.hedis.shared.codec.HedisMessage;
import io.netty.channel.Channel;

import java.util.UUID;

/**
 * HedisSession
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisSession {

    private String id;

    private Channel channel;

    private long createTime;

    public HedisSession(Channel channel) {
        id = UUID.randomUUID().toString();
        createTime = System.currentTimeMillis();
        this.channel = channel;
    }

    public String id() {
        return id;
    }

    public Channel channel() {
        return channel;
    }

    public void writeAndFlush(HedisMessage message) {
        channel.writeAndFlush(message);
    }
}
