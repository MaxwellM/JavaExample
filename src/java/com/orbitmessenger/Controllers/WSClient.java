package com.orbitmessenger.Controllers;

import java.net.URI;
import java.nio.ByteBuffer;
import java.util.ArrayList;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import javafx.scene.layout.VBox;
import org.java_websocket.client.WebSocketClient;
import org.java_websocket.drafts.Draft;
import org.java_websocket.handshake.ServerHandshake;

public class WSClient extends WebSocketClient {

    public String allMessages;

    public WSClient(URI serverUri, Draft draft) {
        super(serverUri, draft);
    }

    public WSClient(URI serverURI) {
        super(serverURI);
    }

    @Override
    public void onOpen(ServerHandshake handshakedata) {
        send("{\"action\":\"getAllMessages\"}");
        //send("Hello, it is me. Mario :)");
        System.out.println("new connection opened");
    }

    @Override
    public void onClose(int code, String reason, boolean remote) {
        System.out.println("closed with exit code " + code + " additional info: " + reason);
    }

    @Override
    public void onMessage(String message) {
        //send("{\"action\":\"getAllMessages\"}");
        //send(formatMessage(message));
        allMessages = message;
        System.out.println("received message: " + message);
    }

    @Override
    public void onMessage(ByteBuffer message) {
        System.out.println("received ByteBuffer");
    }

    @Override
    public void onError(Exception ex) {
        System.err.println("an error occurred:" + ex);
    }

    public String formatMessage(String message) {
        String formattedMessage = "{}";

        return formattedMessage;
    }

    public boolean isAllMessagesEmpty() {
        return allMessages == null;
    }

    public JsonArray getAllMessages() {
        System.out.println("All Messages: " + allMessages);
        JsonParser parser = new JsonParser();
        JsonElement tradeElement = parser.parse(allMessages);
        JsonArray trade = tradeElement.getAsJsonArray();
        System.out.println(trade);
        return trade;
    }

//    public static void main(String[] args) throws URISyntaxException {
//        WebSocketClient client = new WSClient(new URI("ws://localhost:8887"));
//        client.connect();
//    }
}