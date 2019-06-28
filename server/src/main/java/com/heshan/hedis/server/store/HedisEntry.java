package com.heshan.hedis.server.store;

/**
 * HedisEntry
 *
 * @author heshan
 * @date 2019-06-28
 */
public abstract class HedisEntry {

    private long createTime;

    private long expireTime;

    protected HedisEntry() {
        createTime = System.currentTimeMillis();
    }

    public void expire(long expireTime) {
        this.expireTime = expireTime;
    }

    public long expire() {
        return this.expireTime;
    }

    public long createTime() {
        return this.createTime;
    }
}
