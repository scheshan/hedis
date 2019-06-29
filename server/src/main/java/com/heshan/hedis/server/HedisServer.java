package com.heshan.hedis.server;

import com.heshan.hedis.server.handler.CommandHandler;
import com.heshan.hedis.server.handler.SessionHandler;
import com.heshan.hedis.shared.handler.HedisMessageDecoder;
import com.heshan.hedis.shared.handler.HedisMessageEncoder;
import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;

/**
 * HedisServer
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisServer {

    private EventLoopGroup masterGroup;

    private EventLoopGroup workerGroup;

    private Channel channel;

    public HedisServer() {
        masterGroup = new NioEventLoopGroup();
        workerGroup = new NioEventLoopGroup();
    }

    public void start() throws Exception {
        ServerBootstrap bootstrap = new ServerBootstrap()
                .group(masterGroup, workerGroup)
                .childHandler(new ChannelInitializer<Channel>() {
                    @Override
                    protected void initChannel(Channel ch) throws Exception {
                        ChannelPipeline pipeline = ch.pipeline();
                        pipeline.addLast(new HedisMessageDecoder())
                                .addLast(new HedisMessageEncoder())
                                .addLast(new SessionHandler())
                                .addLast(new CommandHandler());
                    }
                })
                .channel(NioServerSocketChannel.class);

        channel = bootstrap.bind(16379).sync().channel();
    }

    public void stop() {
        masterGroup.shutdownGracefully();
        workerGroup.shutdownGracefully();
        channel.close();
    }

    public static void main(String[] args) throws Exception {
        HedisServer server = new HedisServer();
        server.start();

        Runtime.getRuntime().addShutdownHook(server.new ShutdownThread());
    }

    private class ShutdownThread extends Thread {

        @Override
        public void run() {
            HedisServer.this.stop();
        }
    }
}
