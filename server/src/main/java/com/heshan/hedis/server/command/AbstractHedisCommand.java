package com.heshan.hedis.server.command;

import com.heshan.hedis.server.session.SessionManager;

public abstract class AbstractHedisCommand implements HedisCommand {

    protected final CommandFactory commandFactory = CommandFactory.getInstance();

    protected final SessionManager sessionManager = SessionManager.getInstance();
}
