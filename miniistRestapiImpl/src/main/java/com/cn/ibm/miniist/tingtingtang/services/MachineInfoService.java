package com.cn.ibm.miniist.tingtingtang.services;

import com.cn.ibm.miniist.tingtingtang.entities.MachineInfo;
import com.cn.ibm.miniist.tingtingtang.helper.MongoDBHelper;
import com.mongodb.BasicDBObject;
import com.mongodb.BasicDBObjectBuilder;
import com.mongodb.DB;
import com.mongodb.DBCollection;
import com.mongodb.DBObject;

public class MachineInfoService
{
	private static DB db;
	private static MongoDBHelper dbhelper = new MongoDBHelper();
	private static DBCollection collection;
	
//	public void createMachineColl(String coll)
//	{
//		db = dbhelper.getDB();
//		
//        if (db.collectionExists(coll))
//        {
//       	 collection = db.getCollection(coll);
//        }
//        else
//        {
//       	 DBObject options = BasicDBObjectBuilder.start().add("capped", true).add("size", 20001).get();
//       	 collection = db.createCollection(coll, options);
//        	 
//       	 System.out.println("Connect to database successfully");
//        }
//	}
	
	//Insert a new machine into mongodb (need to add boolean to verify the new machine name is available or not)
	public void createNewMachine(MachineInfo machineInfo, String coll)
	{
		//String dbname = "test";
		
		db = dbhelper.getDB();
		if (db.collectionExists(coll))
		{
			collection = db.getCollection(coll);
		}
		else
		{
			DBObject options = BasicDBObjectBuilder.start().add("capped", true).add("size", 20001).get();
	       	collection = db.createCollection(coll, options);
	        System.out.println("Connect to database successfully");
		}
		
		collection.createIndex(new BasicDBObject("name", 1));	
		BasicDBObject doc = new BasicDBObject("name", machineInfo.getMachineName())
				.append("hostName", machineInfo.getHostName())
				.append("ipAddress", machineInfo.getIpAdd())
				.append("description", machineInfo.getDes())
				.append("createTime", machineInfo.getCreateTime())
			    .append("updateTime", machineInfo.getUpdateTime());
				
		collection.insert(doc);		
	}
	
	// Find the required machine info using the unique machine name
	public MachineInfo findMachine(String mName)
	{
		MachineInfo machineInfo = new MachineInfo();
		return machineInfo;
	}
	
	//Update the corresponding machine info using the unique machine name
	public void updateMachine(String mName, MachineInfo machineInfo)
	{
		
	}
	
	//Delete the corresponding machine info in mongodb using the unique machine name
	public void deleteMachine(String mName)
	{
		
	}
}
