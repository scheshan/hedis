package com.heshan.hedis.shared.codec;

/**
 * StringHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class StringHedisMessage implements HedisMessage {

    private String content;

    public StringHedisMessage(String content) {
        this.content = content;
    }

    public String content() {
        return content;
    }
}
