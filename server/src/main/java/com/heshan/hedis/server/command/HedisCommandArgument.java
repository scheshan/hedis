package com.heshan.hedis.server.command;

import com.heshan.hedis.server.session.HedisSession;

/**
 * HedisCommandArgument
 *
 * @author heshan
 * @date 2019-06-30
 */
public class HedisCommandArgument {

    private HedisSession session;

    private String command;

    private String[] args;

    public HedisCommandArgument(HedisSession session, String command, String[] args) {
        this.session = session;
        this.command = command;
        this.args = args;
    }

    public HedisSession session() {
        return session;
    }

    public String command() {
        return command;
    }

    public String[] args() {
        return args;
    }
}
