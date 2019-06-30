package com.heshan.hedis.shared.codec;

/**
 * BatchHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class BatchHedisMessage implements HedisMessage {

    private String content;

    public BatchHedisMessage(String content) {
        this.content = content;
    }

    @Override
    public String content() {
        return content;
    }
}
