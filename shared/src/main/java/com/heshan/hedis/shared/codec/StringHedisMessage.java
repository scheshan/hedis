package com.heshan.hedis.shared.codec;

import com.heshan.hedis.shared.util.HedisMessageUtils;
import io.netty.buffer.ByteBuf;
import io.netty.util.CharsetUtil;

/**
 * StringHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class StringHedisMessage extends AbstractHedisMessage {

    @Override
    protected void doRead(ByteBuf buf) {
        ByteBuf line = HedisMessageUtils.readLine(buf);
        if (line == null) {
            return;
        }

        content = line.readCharSequence(line.readableBytes(), CharsetUtil.UTF_8).toString();
        finish = true;
    }
}
