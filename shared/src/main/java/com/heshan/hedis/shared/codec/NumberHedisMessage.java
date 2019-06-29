package com.heshan.hedis.shared.codec;

/**
 * NumberHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class NumberHedisMessage implements HedisMessage {

    private long value;

    public NumberHedisMessage(long value) {
        this.value = value;
    }

    public long value() {
        return value;
    }
}
