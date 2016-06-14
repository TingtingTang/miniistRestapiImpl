package com.cn.ibm.miniist.tingtingtang.services;

import java.util.HashMap;
import java.util.Map;

import com.cn.ibm.miniist.tingtingtang.entities.MachineInfo;
import com.cn.ibm.miniist.tingtingtang.helper.MongoDBHelper;
import com.mongodb.BasicDBObject;

import com.mongodb.DB;
import com.mongodb.DBCollection;
import com.mongodb.DBCursor;

public class MachineInfoService
{
	private static DB db;
	private static MongoDBHelper dbhelper = new MongoDBHelper();
	
	//Insert a new machine into mongodb (need to add boolean to verify the new machine name is available or not)
	public Boolean createNewMachine(MachineInfo machineInfo, String collName)
	{
		//String dbname = "test";
		
			db = dbhelper.getDB();
			DBCollection collection = db.getCollection(collName);
			
			BasicDBObject findBymName = new BasicDBObject();
			findBymName.put("mName", machineInfo.getMachineName());
			
			DBCursor cursor = collection.find(findBymName);
			if(!cursor.hasNext())
			{
				collection.createIndex(new BasicDBObject("mName", 1));
				BasicDBObject mDoc = new BasicDBObject();
				mDoc.put("mName", machineInfo.getMachineName());
				mDoc.put("hostName", machineInfo.getHostName());
				mDoc.put("ipAddress", machineInfo.getIpAdd());
				mDoc.put("description", machineInfo.getDes());
				mDoc.put("createTime", machineInfo.getCreateTime());
				mDoc.put("updateTime", machineInfo.getUpdateTime());
	
				collection.insert(mDoc);
				//dbhelper.closeMongoSession();
				return true;
			}
			else
			{
				return false;
			}
	}
	
	// Find the required machine info using the unique machine name
	@SuppressWarnings("unchecked")
	public MachineInfo findMachine(String mName, String collName)
	{
		DBCollection collection;
		db = dbhelper.getDB();
		if (db.collectionExists(collName))
		{
			collection = db.getCollection(collName);
			BasicDBObject findBymName = new BasicDBObject();
			findBymName.put("mName", mName);
			
			DBCursor cursor = collection.find(findBymName);
			Map<String, Object> mMap = new HashMap<String, Object>();
			MachineInfo machineInfo = new MachineInfo();
			while(cursor.hasNext())
			{
				mMap = cursor.next().toMap();
			}			
			machineInfo.setMachineName(mName);
			machineInfo.setHostName(mMap.get("hostName").toString());
			machineInfo.setIpAdd(mMap.get("ipAddress").toString());
			machineInfo.setDes(mMap.get("description").toString());
			machineInfo.setCreateTime(mMap.get("createTime").toString());
			machineInfo.setUpdateTime(mMap.get("updateTime").toString());
			machineInfo.setId(mMap.get("_id").toString());
			
			//dbhelper.closeMongoSession();
			
			return machineInfo;
		}
		else
		{
			return null;
		}
	}
	
	//Update the corresponding machine info using the unique machine name
	public Boolean updateMachine(MachineInfo machineInfo, String collName)
	{
		db = dbhelper.getDB();
		DBCollection collection = db.getCollection(collName);
		
		BasicDBObject findBymName = new BasicDBObject();
		findBymName.put("mName", machineInfo.getMachineName());
		
		DBCursor cursor = collection.find(findBymName);
		if(cursor.hasNext())
		{
			BasicDBObject updateFields = new BasicDBObject();
			updateFields.append("hostName", machineInfo.getHostName());
			updateFields.append("ipAddress", machineInfo.getIpAdd());
			updateFields.append("description", machineInfo.getDes());
			updateFields.append("updateTime", machineInfo.getUpdateTime());
			
			BasicDBObject setQuery = new BasicDBObject();
			setQuery.append("$set", updateFields);
			
			collection.update(new BasicDBObject().append("mName", machineInfo.getMachineName()), setQuery, false, true);	
			return true;
		}
		else
		{
			return false;
		}
	}
	
	//Delete the corresponding machine info in mongodb using the unique machine name
	public Boolean deleteMachine(String mName, String collName) 
	{
		db = dbhelper.getDB();
		DBCollection collection = db.getCollection(collName);

		BasicDBObject findBymName = new BasicDBObject();
		findBymName.put("mName", mName);

		DBCursor cursor = collection.find(findBymName);
		if (cursor.hasNext()) 
		{
			BasicDBObject delMachine = new BasicDBObject();
			delMachine.append("mName", mName);

			collection.remove(delMachine);

			return true;
		} 
		else
		{
			return false;
		}
	}
}