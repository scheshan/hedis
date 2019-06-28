package com.heshan.hedis.shared.codec;

import io.netty.buffer.ByteBuf;

/**
 * AbstractHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public abstract class AbstractHedisMessage implements HedisMessage {

    private static final short CRLF = ('\r' << 1) | '\n';

    protected boolean finish;

    protected String content;

    public boolean isFinish() {
        return finish;
    }

    public String content() {
        return content;
    }

    public void read(ByteBuf buf) {
        if (finish) {
            throw new IllegalStateException("Read already completed!");
        }

        doRead(buf);
    }

    protected abstract void doRead(ByteBuf buf);
}
