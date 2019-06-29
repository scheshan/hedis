package com.heshan.hedis.shared.codec;

/**
 * ErrorHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class ErrorHedisMessage implements HedisMessage {

    private String content;

    public ErrorHedisMessage(String content) {
        this.content = content;
    }

    public String content() {
        return content;
    }
}
