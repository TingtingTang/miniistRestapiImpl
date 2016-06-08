package com.cn.ibm.miniist.tingtingtang.helper;

import java.net.UnknownHostException;

import com.mongodb.DB;
import com.mongodb.Mongo;
import com.mongodb.MongoURI;

public class MongoDBHelper 
{
	//private static Mongo mongo;
	private static DB db;
	//private static String host = "localhost";
	private static String host = "198.11.174.68";
	private static int port = 27017;
	//private static String dbname = "test";
	private static String dbname = "test";
	//private static DBCollection collection;
	
	public MongoDBHelper()
	{
		/*
		 * Connect to MongoDB server
		 */
		try 
		{
			//mongo = new Mongo(host, port);
			MongoURI uri = new MongoURI("mongodb://" + host + ":" + port + "/" + dbname);
	    	Mongo mongo = new Mongo(uri);
	    	db = mongo.getDB(uri.getDatabase());
			//db = mongo.getDB(dbname);
		} 
		catch (UnknownHostException e) 
		{
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
//	public DB getDB(String dbname)
//	{
//		return mongo.getDB(dbname);
//	}
	
	public DB getDB()
	{
		return db;
	}
}
