package com.heshan.hedis.shared.codec;

/**
 * NumberHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class NumberHedisMessage implements HedisMessage {

    private long value;

    private String content;

    public NumberHedisMessage(long value) {
        this.value = value;
        content = String.valueOf(value);
    }

    public long value() {
        return value;
    }

    @Override
    public String content() {
        return content;
    }
}
