package com.cn.ibm.tingtingtang.dbconnection.dbConnTest;

import java.sql.Timestamp;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import javax.xml.parsers.DocumentBuilder;

import org.bson.BSONObject;
import org.bson.types.ObjectId;

import com.google.gson.Gson;
import com.mongodb.BasicDBObject;
import com.mongodb.BasicDBObjectBuilder;
import com.mongodb.DB;
import com.mongodb.DBCollection;
import com.mongodb.DBCursor;
import com.mongodb.DBObject;
import com.mongodb.Mongo;
import com.mongodb.MongoClient;
import com.mongodb.MongoClientURI;
import com.mongodb.MongoURI;
import org.bson.types.ObjectId;

public class MongoDBJDBC {
	public static void main(String args[]) {

		try {
			String localhost = "localhost";
			int port = 27017;

			// To connect to mongodb server
			/*
			 * Connect to the local MongoDB and it works well with the below
			 * code; And the version of the local MongoDB is 3.2.4
			 */
			Mongo mongo = new Mongo(localhost, port);
//			MongoURI uri = new MongoURI("mongodb://198.11.174.68:27017/test");
//			Mongo mongo = new Mongo(uri);
			//DB db = mongo.getDB(uri.getDatabase());
			/*
			 * Connect to the remote MongoDB
			 */
			// Mongo mongo = new Mongo("9.125.141.88", 22);

			/*
			 * MongoClientURI uri = new
			 * MongoClientURI("mongodb://tangtt:spr16ing@9.125.141.88:22/test");
			 * MongoClient client = new MongoClient(uri); DB db =
			 * client.getDB(uri.getDatabase());
			 */
			// Now connect to your databases
			DB db = mongo.getDB( "testhub" );
			System.out.println("Connect to database successfully");

//			DBCollection collection;
//			if (db.collectionExists("mlist")) {
//				collection = db.getCollection("mlist");
//				System.out.println("Connect to database successfully");
//			} else {
//				DBObject options = BasicDBObjectBuilder.start().add("capped", true).add("size", 20001).get();
//				collection = db.createCollection("tangtt", options);
//			}
	/*
	 * Create new machine info in the database
	 */
//			DBCollection collection = db.getCollection("mlist");
//			collection.createIndex(new BasicDBObject("mName", 1));
//
//			 BasicDBObject doc = new BasicDBObject("userName", "tangtt")
//					 .append("mName", "w9b2")
//					 .append("hostName", "www.w9b2.com")
//					 .append("ipAddress", "127.0.0.2")
//					 .append("description", "w9b2")
//					 .append("createTime", "20160605")
//					 .append("updateTime", "20160607");
//			 collection.insert(doc);
//			System.out.println("Connect to database successfully");
			
			
			/*
			 * find an object by machine name in a collection
			 */
//			DBCollection collection;
//
//			collection = db.getCollection("mlist");
//			BasicDBObject findBymName = new BasicDBObject();
//			findBymName.put("mName", "w9b1");
//
//			DBCursor cursor = collection.find(findBymName);
//			Map<String, Object> mMap = new HashMap<String, Object>();
//			MachineInfo machineInfo = new MachineInfo();
//			while (cursor.hasNext()) {
//				mMap = cursor.next().toMap();
//			}
//			machineInfo.setMachineName("w9b1");
//			machineInfo.setUserName(mMap.get("userName").toString());
//			machineInfo.setHostName(mMap.get("hostName").toString());
//			machineInfo.setIpAdd(mMap.get("ipAddress").toString());
//			machineInfo.setDes(mMap.get("description").toString());
//			machineInfo.setCreateTime(mMap.get("createTime").toString());
//			machineInfo.setUpdateTime(mMap.get("updateTime").toString());
//			machineInfo.setId(mMap.get("_id").toString());
//			System.out.println("Connect to database successfully !");
//			System.out.println(machineInfo.getUserName());
			
			/*
			 * find all machine information under specific user
			 */

			DBCollection collection = db.getCollection("mlist");
			DBObject query = new BasicDBObject("userName","tangtt");
			DBCursor cursor = collection.find(query);		
			
			List<MachineInfo> list = new ArrayList<MachineInfo>();

			while(cursor.hasNext())
			{
				MachineInfo machineInfo = new MachineInfo();
				DBObject o = cursor.next();
				
				machineInfo.setMachineName((String) o.get("mName"));
				machineInfo.setHostName((String) o.get("hostName"));
				machineInfo.setIpAdd((String) o.get("ipAddress"));
				machineInfo.setDes((String) o.get("description"));
				machineInfo.setCreateTime((String) o.get("createTime"));
				machineInfo.setUpdateTime((String) o.get("updateTime"));
				
				ObjectId _id = (ObjectId) o.get("_id");
				machineInfo.setId(_id.toString());
				
				machineInfo.setUserName("tangtt");
				
				System.out.println(machineInfo.getId());
				System.out.println(machineInfo.toString());
				
				list.add(machineInfo);
			}
			
			cursor.close();
			
		/*
		 * update machine information of user	
		 */
//			DBCollection collection = db.getCollection("mlist");
//			
//			BasicDBObject findBymName = new BasicDBObject();
//			findBymName.put("mName", "w9b1");
//			
//			DBCursor cursor = collection.find(findBymName);
//			if(cursor.hasNext())
//			{
//				BasicDBObject updateFields = new BasicDBObject();
//				updateFields.append("hostName", "www.w9b1.com");
//				updateFields.append("ipAddress", "127.0.0.1");
//				updateFields.append("description", "w9b1");
//				updateFields.append("updateTime", "20160712");
//				
//				BasicDBObject setQuery = new BasicDBObject();
//				setQuery.append("$set", updateFields);
//				
//				collection.update(new BasicDBObject().append("mName", "w9b1"), setQuery, false, true);
//				//collection.update(new BasicDBObject().append("userName", "tangtt"), setQuery, false, true);
//				System.out.println("Connect to database successfully !");
//			}
			
			
	/*
	 * Delete document in MongoDB
	 */
//			DBCollection collection = db.getCollection("mlist");
//			BasicDBObject findBymName = new BasicDBObject();
//			findBymName.put("mName", "w9b1");
//
//			DBCursor cursor = collection.find(findBymName);
//			if (cursor.hasNext()) 
//			{
//				BasicDBObject delMachine = new BasicDBObject();
//				delMachine.append("mName", "w9b1");
//
//				collection.remove(delMachine);
//				System.out.println("Connect to database successfully !");
//			} 
		}
			catch (Exception e)
		{
			System.err.println(e.getClass().getName() + ": " + e.getMessage());
		}
	}
}