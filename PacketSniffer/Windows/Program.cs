using System;
using System.Collections;
using System.Collections.Generic;
using System.Collections.Concurrent;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Collections.Concurrent;
using System.Net;
using PcapDotNet.Core;
using PcapDotNet.Packets;
using PcapDotNet.Packets.IpV4;
using PcapDotNet.Packets.Transport;
using PcapDotNet.Packets.Http;
using System.Text;
using RabbitMQ.Client;
using Newtonsoft.Json;

namespace Guardian
{
    public class NetworkListener
    {
        public void Listen(BlockingCollection<GuardianRequest> col)
        {
            var writer = new Producer<GuardianRequest>(col);
            var reader = new Consumer<GuardianRequest>(col);

            Task.Factory.StartNew(() =>
            {
                GuardianEntity guardian = new GuardianEntity();

                while (true)
                {
                    GuardianRequest request = reader.Consume();
                    guardian.Add(request);

<<<<<<< HEAD
                    Console.WriteLine("Connecting");
                    try
                    {
                        ConnectionFactory factory = new ConnectionFactory();

                        factory.Uri = "amqp://guest:guest@localhost:5672/";

                        using (IConnection conn = factory.CreateConnection())
                        {
                            IModel channel = conn.CreateModel();
                            Console.WriteLine("Publishing message");
                            channel.BasicPublish("test-exchange", "test-key", null, Encoding.Default.GetBytes(JsonConvert.SerializeObject(request)));
                        }
=======
                    var requests = guardian.GetRequestsToday();
                    foreach (var r in requests)
                    {
                        Console.WriteLine("Host: {0} started: {1} ended: {2} totalvisits: {3}", r.Host, r.Started, r.Ended, r.Count);
>>>>>>> 3952e23c477caa98e6a58d87c208e780db46350f
                    }
                    catch (Exception e)
                    {
                        Console.WriteLine(e.ToString());
                    }
                }
            });

            IList<LivePacketDevice> allDevices = LivePacketDevice.AllLocalMachine;
            LivePacketDevice device = allDevices[2];

            using (PacketCommunicator communicator = device.Open(65536, PacketDeviceOpenAttributes.Promiscuous, 1000))
            {
                Console.WriteLine("Listening on " + device.Description);


                communicator.SetFilter("ip and tcp");


                communicator.ReceivePackets(0, ((Packet packet) =>
                {
                    try
                    {
                        if (packet != null)
                        {
                            IpV4Datagram ip = packet.Ethernet.IpV4;
                            if (new[] { 80 }.Contains(ip.Tcp.DestinationPort))
                            {
                                TcpDatagram tcpDatagram = ip.Tcp;
                                HttpDatagram httpDatagram = tcpDatagram.Http;


                                if (httpDatagram.IsRequest)
                                {
                                    HttpRequestDatagram httpRequestDatagram = (HttpRequestDatagram)httpDatagram;

                                    if (httpRequestDatagram.Uri == "/" && httpRequestDatagram.IsValid)
                                    {
                                        var host = httpRequestDatagram.Header["Host"];

                                        GuardianRequest req = new GuardianRequest
                                        {
                                            Port = tcpDatagram.DestinationPort,
                                            Host = host.ValueString,
                                            TimeStamp = packet.Timestamp
                                        };

                                        writer.Produce(req);

                                    }
                                }

                                if (httpDatagram.IsResponse)
                                {
                                    HttpResponseDatagram httpRequestDatagram = (HttpResponseDatagram)httpDatagram;
                                    throw new Exception("asd");
                                }
                            }




                        }
                    }
                    catch (Exception e) { }
                }));
            }
        }
    }

    class Program
    {
        static void Main(string[] args)
        {
<<<<<<< HEAD

            NetworkListener listener = new NetworkListener();
            var request = new BlockingCollection<GuardianRequest>();
            listener.Listen(request);
=======
            var listener = new NetworkListener();
            listener.Listen(new System.Collections.Concurrent.BlockingCollection<GuardianRequest>());

>>>>>>> 3952e23c477caa98e6a58d87c208e780db46350f
        }

        // Print all the available information on the given interface
        private static void DevicePrint(IPacketDevice device)
        {
            // Name
            Console.WriteLine(device.Name);

            // Description
            if (device.Description != null)
                Console.WriteLine("\tDescription: " + device.Description);

            // Loopback Address
            Console.WriteLine("\tLoopback: " +
                              (((device.Attributes & DeviceAttributes.Loopback) == DeviceAttributes.Loopback)
                                   ? "yes"
                                   : "no"));

            // IP addresses
            foreach (DeviceAddress address in device.Addresses)
            {
                Console.WriteLine("\tAddress Family: " + address.Address.Family);

                if (address.Address != null)
                    Console.WriteLine(("\tAddress: " + address.Address));
                if (address.Netmask != null)
                    Console.WriteLine(("\tNetmask: " + address.Netmask));
                if (address.Broadcast != null)
                    Console.WriteLine(("\tBroadcast Address: " + address.Broadcast));
                if (address.Destination != null)
                    Console.WriteLine(("\tDestination Address: " + address.Destination));
            }
            Console.WriteLine();
        }







    }
}
