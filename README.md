### Inspired by [dtm](https://dtm.pub/ref/gozero.html), Saga Distributed Transactions Based on Redis Distributed Locks, Taking go-zero as an Example

github：https://github.com/H-H1/SAGA-go-zero

English | [简体中文](README-cn.md)

##### Under a pressure test of 4000 - 5000 TPS, 98% of eventual consistency can be guaranteed.

> In the common CAP theory, C stands for Consistency, A for Availability, and P for Partition Tolerance. Generally, P and A are ensured while sacrificing C and opting for eventual consistency. This article achieves AP and eventual consistency and incorporates a retry mechanism for compensation operations.
>I previously studied DTM but discovered an issue. In high-concurrency scenarios, DTM transactions are unstable, with numerous compensation operations failing, which has been verified multiple times. Hence, this article was written. The following is the sequence diagram.
> 
>![image-20241214202609524](https://github.com/user-attachments/assets/cd42f34d-a3d8-4f31-823e-a12eaddfa56e)

- Here, I prepare a DTM pressure test with a concurrency of 1000 per second.
- The stock inventory is 100. If the execution is successful, the inventory will decrease by 1. During rollback compensation, it will increase by 1.

![image-20241214200601723](https://github.com/user-attachments/assets/90856cd0-8150-4a9a-aa9a-1ca9bb7148cd)


- If the insertion in the order is successful, a new record will be created, and the value of row_state will be 0. During rollback compensation, row_state will be -1.

![image-20241214200354019](https://github.com/user-attachments/assets/742dbe3d-e76c-4ce1-a7d3-96ebde8f2d09)


- Looking at the pressure test results, 217 were successful.

![image-20241214200534660](https://github.com/user-attachments/assets/b6999b00-6ac4-4037-9f3a-2efe48a30e72)


- There is no problem with the stock. The deduction was completed successfully.

![image-20241214200706228](https://github.com/user-attachments/assets/46e99d56-f299-4606-8096-0016a7c74ae5)


- However, there is a problem with the order. A large number of records in the successful state were inserted, and many compensations failed.

![image-20241214201343030](https://github.com/user-attachments/assets/4a8dd73a-72c7-4416-b17b-78c308278b65)


#### However, there is a problem with the order. A large number of records in the successful state were inserted, and many compensations failed.

This transaction can currently provide consistency for 2 microservices. For more than 2 services, I think they can be merged first. The specific sequence diagram is as follows.

![image-20241214202609524](https://github.com/user-attachments/assets/cd42f34d-a3d8-4f31-823e-a12eaddfa56e)

1. Two errors are used to represent the status of the confirmation microservices respectively.

   

   1. If Create(l.ctx, createOrderReq) fails, the service directly returns failure, and nothing happens.
   2. If Create(l.ctx, createOrderReq) succeeds with row_state being 0, but l.svcCtx.StockRpc.Deduct(l.ctx, deductReq) fails, compensation for the order will be performed, and the compensation includes retries.
   3. Both microservices incorporate the method of distributed lock + uid to prevent misdeletion caused by latency.

- Performing a pressure test, it can be seen that due to the existence of the lock, the number of successes is only 4.

![image-20241214204252365](https://github.com/user-attachments/assets/f2de42e3-ce75-468f-ace1-5d73c80c8b38)


However, there will be no inconsistency in the order and stock caused by compensation failure. 4 were deducted, and all are in state 0.

![image-20241214204419091](https://github.com/user-attachments/assets/bb53d5ea-09ce-4a54-aca5-0f694326c450)

![image-20241214204433122](https://github.com/user-attachments/assets/1909edd3-0ee3-4282-b652-391618e014e1)



The results of multiple pressure tests are consistent.

![image-20241214204604732](https://github.com/user-attachments/assets/975df20c-76b4-42da-b80a-0eb13b881130)

![image-20241214204613884](https://github.com/user-attachments/assets/e748435e-fc9a-488f-95c0-176518141234)


![image-20241214204626565](https://github.com/user-attachments/assets/2046eedb-c117-41f5-8720-07246d1c0946)



Under a pressure test of 4000 - 5000 TPS, 98% of eventual consistency can be guaranteed.

![image-20241214210804031](https://github.com/user-attachments/assets/f6485935-cff3-481c-b66a-7b2ba6aaebb1)


![image-20241214210811660](https://github.com/user-attachments/assets/ab54f6f0-331b-4fa9-b441-a2131b722b6d)

![image-20241214210819236](https://github.com/user-attachments/assets/78009ccb-2af8-4793-9243-9bf21cf0b445)

- Acknowledgments

1. [dtm](https://github.com/dtm-labs/dtm)

2. [go-zero](https://github.com/zeromicro/go-zero)

3. [looklook](https://github.com/zeromicro/go-zero)

   

