package com.heshan.hedis.server.command;

import com.heshan.hedis.server.session.HedisSession;

public interface HedisCommand {

    void execute(HedisSession session, String[] args);
}
