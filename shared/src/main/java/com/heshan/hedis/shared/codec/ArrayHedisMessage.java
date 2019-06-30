package com.heshan.hedis.shared.codec;

import java.util.Collection;

/**
 * ArrayHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class ArrayHedisMessage implements HedisMessage {

    private Collection<HedisMessage> messages;

    public ArrayHedisMessage(Collection<HedisMessage> messages) {
        this.messages = messages;
    }

    public int size() {
        return messages.size();
    }

    public Collection<HedisMessage> messages() {
        return messages;
    }

    @Override
    public String content() {
        throw new UnsupportedOperationException();
    }
}
