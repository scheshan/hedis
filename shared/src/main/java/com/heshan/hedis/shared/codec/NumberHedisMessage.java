package com.heshan.hedis.shared.codec;

import com.heshan.hedis.shared.util.HedisMessageUtils;
import io.netty.buffer.ByteBuf;

/**
 * NumberHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class NumberHedisMessage extends AbstractHedisMessage {

    private long value;

    @Override
    protected void doRead(ByteBuf buf) {
        ByteBuf line = HedisMessageUtils.readLine(buf);
        if (line == null) {
            return;
        }

        value = HedisMessageUtils.readNumber(line);
        content = String.valueOf(value);
        finish = true;
    }
}
