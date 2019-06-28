package com.heshan.hedis.shared.codec;

import com.heshan.hedis.shared.exception.HedisProtocolException;
import com.heshan.hedis.shared.util.HedisMessageUtils;
import io.netty.buffer.ByteBuf;
import io.netty.util.CharsetUtil;

/**
 * BatchHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class BatchHedisMessage extends AbstractHedisMessage {

    private int length = -2;

    @Override
    protected void doRead(ByteBuf buf) {
        if (length < -1) {
            ByteBuf line = HedisMessageUtils.readLine(buf);
            if (line == null) {
                return;
            }

            length = (int) HedisMessageUtils.readNumber(line);
            if (length < -1) {
                throw new HedisProtocolException();
            }
        }

        if (length == -1) {
            finish = true;
            return;
        }

        ByteBuf data = HedisMessageUtils.readLength(buf, length);
        if (data == null) {
            return;
        }

        content = data.readCharSequence(data.readableBytes(), CharsetUtil.UTF_8).toString();
        finish = true;
    }
}
