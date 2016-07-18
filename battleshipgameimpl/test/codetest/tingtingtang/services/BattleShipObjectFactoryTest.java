/**
 * 
 */
package codetest.tingtingtang.services;

import static org.junit.Assert.*;

import org.junit.BeforeClass;
import org.junit.Test;

import codetest.tingtingtang.entities.ShipInfo;
import codetest.tingtingtang.enums.ShipType;

/**
 * @author TingTing Tang
 *
 */
public class BattleShipObjectFactoryTest 
{

	/**
	 * @throws java.lang.Exception
	 */
	@BeforeClass
	public static void setUpBeforeClass() throws Exception 
	{
	}

	/**
	 * Test method for {@link codetest.tingtingtang.services.BattleShipObjectFactory#generateShip(codetest.tingtingtang.enums.ShipType)}.
	 */
	@Test
	public void testGenerateShip() 
	{
		ShipInfo shipInfoB = BattleShipObjectFactory.generateShip(ShipType.BATTLESHIP);
		assertNotNull(shipInfoB);
		assertTrue(shipInfoB.getShipLength() == 5);
		
		ShipInfo shipInfoD = BattleShipObjectFactory.generateShip(ShipType.DESTROYER);
		assertNotNull(shipInfoD);
		assertTrue(shipInfoD.getShipLength() == 4);
	}

	/**
	 * Test method for {@link codetest.tingtingtang.services.BattleShipObjectFactory#generateBoard(int, int)}.
	 */
	@Test
	public void testGenerateBoard() 
	{
		boolean[][] board = BattleShipObjectFactory.generateBoard(10, 10);
		assertNotNull(board);
		assertTrue(board.length == 10);
	}

}
