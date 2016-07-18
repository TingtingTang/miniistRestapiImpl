package codetest.tingtingtang.services;

import static org.junit.Assert.*;

import java.util.ArrayList;
import java.util.List;

import org.junit.BeforeClass;
import org.junit.Test;

import codetest.tingtingtang.entities.ShipInfo;
import codetest.tingtingtang.enums.RowOfBoard;
import codetest.tingtingtang.enums.ShipType;
import codetest.tingtingtang.enums.TargetOutcome;

public class BattleShipServiceTest {

	static ShipInfo shipInfoB = null;
	static ShipInfo shipInfoDO = null;
	static ShipInfo shipInfoDT = null;
	static boolean[][] board = null;
	@BeforeClass
	public static void setUpBeforeClass() throws Exception 
	{
		shipInfoB = BattleShipObjectFactory.generateShip(ShipType.BATTLESHIP);
		shipInfoDO = BattleShipObjectFactory.generateShip(ShipType.DESTROYER);
		shipInfoDT = BattleShipObjectFactory.generateShip(ShipType.DESTROYER);
		board = BattleShipObjectFactory.generateBoard(10, 10);
	}

	@Test
	public void testBattleShipService() 
	{
		List<ShipInfo> ships = new ArrayList<ShipInfo>();
		ships.add(shipInfoB);
		ships.add(shipInfoDO);
		ships.add(shipInfoDT);
		
		BattleShipService battleShipService = new BattleShipService(ships, board);
		assertNotNull(battleShipService.board);
		assertNotNull(battleShipService.ships);
		
	}

	@Test
	public void testAnalyseTarget() 
	{
		List<ShipInfo> ships = new ArrayList<ShipInfo>();
		ships.add(shipInfoB);
		ships.add(shipInfoDO);
		ships.add(shipInfoDT);
		
		BattleShipService battleShipService = new BattleShipService(ships, board);
		String outcome = battleShipService.analyseTarget(getStringFromValue(shipInfoB.getStartingPointX())+(shipInfoB.getStartingPointY()+1));
		assertTrue(outcome.equalsIgnoreCase(TargetOutcome.HIT.name()));
	}
	
	private String getStringFromValue(int row)
	{
		if(row == RowOfBoard.A.getValue())
			return RowOfBoard.A.name();
		if(row == RowOfBoard.B.getValue())
			return RowOfBoard.B.name();
		if(row == RowOfBoard.C.getValue())
			return RowOfBoard.C.name();
		if(row == RowOfBoard.D.getValue())
			return RowOfBoard.D.name();
		if(row == RowOfBoard.E.getValue())
			return RowOfBoard.E.name();
		if(row == RowOfBoard.F.getValue())
			return RowOfBoard.F.name();
		if(row == RowOfBoard.G.getValue())
			return RowOfBoard.G.name();
		if(row == RowOfBoard.H.getValue())
			return RowOfBoard.H.name();
		if(row == RowOfBoard.I.getValue())
			return RowOfBoard.I.name();
		if(row == RowOfBoard.J.getValue())
			return RowOfBoard.J.name();
		return null;
	}

}
