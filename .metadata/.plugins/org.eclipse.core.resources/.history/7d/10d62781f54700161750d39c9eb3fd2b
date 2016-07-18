package com.cn.ibm.miniist.tingtingtang.helper;

import java.net.UnknownHostException;

import com.mongodb.DB;
import com.mongodb.Mongo;
import com.mongodb.MongoURI;

public class MongoDBHelper 
{
	//private static Mongo mongo;
	private static DB db;
	/*
	 * This host is used to connect to DB locally
	 */
	private static String host = "localhost";
	/*
	 * This host is used to connect to DB remotely
	 */
	//private static String host = "9.125.141.88";
	private static int port = 27017;
	//private static String dbname = "test";
	private static String dbname = "testhub";
	//private static DBCollection collection;
	private static Mongo mongo;
	
	public MongoDBHelper()
	{
		/*
		 * Connect to MongoDB server
		 */
		try 
		{
			/*
			 * This is used for connect to mongodb remotely
			 */
//			mongo = new Mongo(host, port);
//			MongoURI uri = new MongoURI("mongodb://" + host + ":" + port + "/" + dbname);
//	    	mongo = new Mongo(uri);
//	    	db = mongo.getDB(uri.getDatabase());
			
			/*
			 * This is used for connect to mongodb locally
			 */
			mongo = new Mongo(host,port);
			db = mongo.getDB(dbname);
		} 
		catch (UnknownHostException e) 
		{
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
	public DB getDB()
	{
		return db;
	}
	
	public void closeMongoSession()
	{
		mongo.close();
	}
}
